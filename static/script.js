const dropzone = document.getElementById('dropzone');
const fileInput = document.getElementById('fileInput');
const browseBtn = document.getElementById('browseBtn');
const progress = document.getElementById('progress');
const progressBar = document.getElementById('progressBar');
const result = document.getElementById('result');
const error = document.getElementById('error');

// Открытие выбора файла
browseBtn.addEventListener('click', (e) => {
    e.stopPropagation();
    fileInput.click();
});

dropzone.addEventListener('click', () => fileInput.click());

fileInput.addEventListener('change', () => {
    if (fileInput.files.length > 0) {
        uploadFile(fileInput.files[0]);
    }
});

// Drag & Drop
['dragenter', 'dragover'].forEach(evt =>
    dropzone.addEventListener(evt, (e) => {
        e.preventDefault();
        dropzone.classList.add('dragover');
    })
);

['dragleave', 'drop'].forEach(evt =>
    dropzone.addEventListener(evt, (e) => {
        e.preventDefault();
        dropzone.classList.remove('dragover');
    })
);

dropzone.addEventListener('drop', (e) => {
    if (e.dataTransfer.files.length > 0) {
        uploadFile(e.dataTransfer.files[0]);
    }
});

function uploadFile(file) {
    // Сброс UI
    result.classList.add('hidden');
    error.classList.add('hidden');
    progress.classList.remove('hidden');
    progressBar.style.width = '0%';

    const maxSize = 500 * 1024 * 1024;
    if (file.size > maxSize) {
        showError('Файл слишком большой (максимум 500 МБ)');
        return;
    }

    const formData = new FormData();
    formData.append('file', file);
    formData.append('originalName', file.name);

    const xhr = new XMLHttpRequest();
    xhr.open('POST', '/api/upload');

    xhr.upload.addEventListener('progress', (e) => {
        if (e.lengthComputable) {
            const percent = (e.loaded / e.total) * 100;
            progressBar.style.width = percent + '%';
        }
    });

    xhr.addEventListener('load', () => {
        progress.classList.add('hidden');
        if (xhr.status === 200) {
            try {
                const data = JSON.parse(xhr.responseText);
                showResult(data);
            } catch {
                showError('Некорректный ответ сервера');
            }
        } else {
            showError(xhr.responseText || 'Ошибка загрузки');
        }
    });

    xhr.addEventListener('error', () => {
        progress.classList.add('hidden');
        showError('Ошибка сети');
    });

    xhr.send(formData);
}

function showResult(data) {
    const fullUrl = window.location.origin + data.url;
    const minutes = Math.floor(data.expires_in / 60);
    const durationText = minutes >= 60
    ? `${Math.floor(minutes / 60)} ч ${minutes % 60 ? (minutes % 60) + ' мин' : ''}`.trim()
    : `${minutes} мин`;

    result.innerHTML = `
        <h3>✅ Файл загружен</h3>
        <div class="result-row">
            <input type="text" id="shareUrl" value="${fullUrl}" readonly>
            <button id="copyBtn">Копировать</button>
        </div>
        <div class="meta">
            <strong>${escapeHtml(data.name)}</strong> — ${formatSize(data.size)}<br>
            Ссылка будет активна ${durationText}.
        </div>
    `;
    result.classList.remove('hidden');

    document.getElementById('copyBtn').addEventListener('click', () => {
        const input = document.getElementById('shareUrl');
        input.select();
        navigator.clipboard.writeText(input.value).then(() => {
            const btn = document.getElementById('copyBtn');
            btn.textContent = 'Скопировано!';
            setTimeout(() => (btn.textContent = 'Копировать'), 1500);
        });
    });
}

function showError(msg) {
    error.textContent = msg;
    error.classList.remove('hidden');
}

function formatSize(bytes) {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
}

function escapeHtml(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}
