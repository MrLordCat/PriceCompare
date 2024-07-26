/* scripts.js */

document.addEventListener('DOMContentLoaded', function() {
    const priceAdjustment = 2.939; // Процентная разница для корректировки

    const amazonPriceCells = document.querySelectorAll('td.price-amazon');

    amazonPriceCells.forEach(cell => {
        const originalPrice = parseFloat(cell.innerText);

        if (!isNaN(originalPrice)) {
            const adjustedPrice = originalPrice + (originalPrice * (priceAdjustment / 100));
            cell.innerHTML = adjustedPrice.toFixed(2) + ` <span style="font-size: smaller; color: gray;">(${originalPrice.toFixed(2)})</span>`; // Обновляем цену с учетом процентной разницы и добавляем исходную цену в скобках
        }
    });
});


document.addEventListener('DOMContentLoaded', function() {
    const rows = document.querySelectorAll('table tbody tr');
    rows.forEach(row => {
        const priceCell = row.querySelector('td.price');
        const priceAmazonCell = row.querySelector('td.price-amazon');
        const priceDiffCell = row.querySelector('td.price-diff');

        if (priceCell && priceAmazonCell && priceDiffCell) {
            const price = parseFloat(priceCell.innerText);
            const priceAmazon = parseFloat(priceAmazonCell.innerText);

            if (isNaN(priceAmazon) || priceAmazon == 0) {
                priceDiffCell.innerText = 'N/A';
                priceAmazonCell.innerText = 'N/A';
            } else {
                const priceDiff = price - priceAmazon;
                priceDiffCell.innerText = priceDiff.toFixed(2);

                if (priceDiff > 0) {
                    if (priceDiff <= 10) {
                        priceDiffCell.classList.add('positive-1');
                    } else if (priceDiff <= 30) {
                        priceDiffCell.classList.add('positive-2');
                    } else {
                        priceDiffCell.classList.add('positive-3');
                    }
                } else {
                    priceDiffCell.classList.add('negative');
                }
            }
        }
    });
});
document.addEventListener('DOMContentLoaded', function() {
    const rows = document.querySelectorAll('table tbody tr');
    rows.forEach(row => {
        const usedCell = row.querySelector('td.used');

        if (usedCell) {
            if (usedCell.innerText.toLowerCase() === 'yes') {
                usedCell.classList.add('used-yes');
            } else {
                usedCell.classList.add('used-no');
            }
        }
    });
});
document.addEventListener('DOMContentLoaded', function() {
    const amazonCells = document.querySelectorAll('td.link-amazon');

    amazonCells.forEach(cell => {
        if (cell.innerText === 'N/A') {
            cell.classList.add('link-amazon-no');
        } else {
            cell.classList.add('link-amazon-yes');
        }
    });
});
document.addEventListener('DOMContentLoaded', function() {
    const rows = document.querySelectorAll('table tbody tr');

    rows.forEach(row => {
        const deliveryTimeCell = row.querySelector('td.delivery-time');
        const titleCell = row.querySelector('td.title-amazon');

        if (deliveryTimeCell && titleCell) {
            if (deliveryTimeCell.innerText === 'N/A') {
                titleCell.classList.add('title-amazon-no');
            } else {
                titleCell.classList.add('title-amazon-yes');
            }
        }
    });
});


document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('form.update-amazon-link').forEach(form => {
        form.addEventListener('submit', function(event) {
            event.preventDefault();

            const formData = new FormData(form);
            const title = formData.get('title');
            const category = formData.get('category');
            const linkAmazon = formData.get('linkAmazon');
            
            fetch('/update', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.text();
            })
            .then(data => {
                console.log(`Updated Amazon link for title: ${title}`);
                form.querySelector('input[name="linkAmazon"]').value = ''; // Clear the input
            })
            .catch(error => {
                console.error('There was a problem with the fetch operation:', error);
            });
        });
    });
});


document.querySelectorAll('form.update-fb-link').forEach(form => {
    form.addEventListener('submit', function(event) {
        event.preventDefault();

        const formData = new FormData(form);
        const title = formData.get('title');
        const category = formData.get('category');
        const fbLink = formData.get('FBLink');
        
        fetch('/update-fb', {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.text();
        })
        .then(data => {
            console.log(`Updated FB link for title: ${title}`);
            form.querySelector('input[name="FBLink"]').value = ''; // Clear the input
        })
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
        });
    });
});

document.addEventListener("DOMContentLoaded", function() {
    var rows = document.querySelectorAll("table tbody tr");
    rows.forEach(function(row) {
        var priceElement = row.querySelector("td.price");
        if (priceElement) {
            var price = parseFloat(priceElement.innerText);
            var priceMinus15 = price * 0.85;
            var priceMinus15Element = row.querySelector("td.price-minus-15");
            if (priceMinus15Element) {
                priceMinus15Element.innerText = priceMinus15.toFixed(2);
            }
        }
    });
});
