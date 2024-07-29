/* scripts.js */


document.addEventListener('DOMContentLoaded', function() {
    const priceElements = document.querySelectorAll('.price');
    priceElements.forEach(priceElement => {
        const price = parseFloat(priceElement.getAttribute('data-price'));
        if (!isNaN(price)) {
            const priceMinus15 = (price * 0.85).toFixed(2);
            const priceMinus15Element = priceElement.closest('tr').querySelector('[data-price-minus-15]');
            if (priceMinus15Element) {
                priceMinus15Element.textContent = priceMinus15 + '€';
            }
            // Добавление значка евро к исходной цене
            priceElement.textContent += '€';
        }
    });
});

document.getElementById("update-amazon-form").addEventListener("submit", function (e) {
    e.preventDefault();

    var form = this;
    var formData = new FormData(form);
    var category = formData.get("category");

    fetch(form.action, {
        method: form.method,
        body: formData
    }).then(response => {
        if (response.ok) {
            checkProgress(category);
        }
    }).catch(error => {
        console.error("Error:", error);
    });
});

document.getElementById("update-amazon-form").addEventListener("submit", function (e) {
    e.preventDefault();

    var form = this;
    var formData = new FormData(form);
    var category = formData.get("category");

    fetch(form.action, {
        method: form.method,
        body: formData
    }).then(response => {
        if (response.ok) {
            // Показать прогресс-бар при начале загрузки
            document.getElementById("progress-container").style.display = "block";
            checkProgress(category);
        }
    }).catch(error => {
        console.error("Error:", error);
    });
});

function checkProgress(category) {
    var progressBar = document.getElementById("progress-bar");

    function updateProgress() {
        fetch("/get-progress?category=" + category)
            .then(response => response.json())
            .then(data => {
                var progress = data.progress;
                progressBar.style.width = progress + "%";
                progressBar.textContent = progress + "%";

                if (progress < 100) {
                    setTimeout(updateProgress, 1000);
                } else {
                    // Обновление страницы после завершения загрузки
                    location.reload();
                }
            })
            .catch(error => {
                console.error("Error:", error);
            });
    }

    updateProgress();
}
