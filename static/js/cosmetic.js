document.addEventListener("DOMContentLoaded", function() {
    // Получение всех строк таблицы
    const rows = document.querySelectorAll('tr');

    rows.forEach(row => {
        // Получение первой цены
        const priceElement = row.querySelector('.price');
        if (!priceElement) return;
        const price = parseFloat(priceElement.getAttribute('data-price'));

        // Получение второй цены
        const amazonPriceElement = row.querySelector('.price-amazon');
        if (!amazonPriceElement) return;
        let amazonPrice = amazonPriceElement.textContent.trim();
        amazonPrice = (amazonPrice !== 'N/A' && parseFloat(amazonPrice) !== 0) ? parseFloat(amazonPrice) : null;

        // Элемент для записи разницы цен
        const priceDiffElement = row.querySelector('.price-diff');
        if (!priceDiffElement) return;

        // Вычисление разницы и запись результата
        if (amazonPrice !== null) {
            const priceDifference = price - amazonPrice;
            priceDiffElement.textContent = priceDifference + '€';
        } else {
            priceDiffElement.textContent = 'N/A';
        }

        // Изменение цвета в зависимости от значения разницы
        var value = priceDiffElement.textContent.trim();
        var number = parseFloat(value);

        if (!isNaN(number)) {
            var red = 0, green = 0;
            if (number < 0) {
                red = Math.min(255, Math.max(0, Math.floor(255 * (1 - (Math.abs(number) / 100)))));
                priceDiffElement.style.backgroundColor = `rgba(${red}, 0, 10, 0.7)`; // Добавлена прозрачность 0.7
                priceDiffElement.style.fontWeight = 'bold';
            } else if (number > 0) {
                green = Math.min(255, Math.max(0, Math.floor(255 * (1 - (number / 100)))));
                priceDiffElement.style.backgroundColor = `rgba(0, ${green}, 0, 0.7)`; // Добавлена прозрачность 0.7
                priceDiffElement.style.fontWeight = 'bold';
            } else {
                priceDiffElement.style.backgroundColor = 'rgba(0, 0, 0, 0.7)'; // Добавлена прозрачность 0.7
                priceDiffElement.style.fontWeight = 'bold';
            }
            
            
            
        }
    });
});

document.addEventListener("DOMContentLoaded", function() {
    document.querySelectorAll('.used').forEach(function(element) {
        if (element.textContent.trim() === "Yes") {
            element.style.backgroundColor = 'rgba(174, 32, 32, 0.676)';
        } else {
            element.style.backgroundColor = 'rgba(33, 119, 33, 0.823)';
        }
    });
});

document.addEventListener("DOMContentLoaded", function() {
    document.querySelectorAll('tr').forEach(function(row) {
        // Проверка доступности на Amazon и использованности
        const statusElement = row.querySelector('.status-available, .status-unavailable');
        const usedElement = row.querySelector('.used');
        const linkElement = row.querySelector('.link-na, .link-amazon');
        const titleAmazonElement = row.querySelector('.title-amazon');

        // Если нет titleAmazonElement, пропускаем эту строку
        if (!titleAmazonElement) return;

        const statusText = statusElement ? statusElement.textContent.trim() : '';
        const usedText = usedElement ? usedElement.textContent.trim() : '';
        const linkText = linkElement ? linkElement.textContent.trim() : '';

        // Проверка условий и установка фона
        if (statusText === "Unavailable" || usedText === "Yes" || linkText === "N/A") {
            titleAmazonElement.style.backgroundColor = 'rgba(174, 32, 32, 0.676)';
        } else {
            titleAmazonElement.style.backgroundColor = 'rgba(33, 119, 33, 0.823)';
        }
    });
});
