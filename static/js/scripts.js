/* scripts.js */


document.addEventListener('DOMContentLoaded', function() {
    const priceElements = document.querySelectorAll('.price');
    priceElements.forEach(priceElement => {
        const price = parseFloat(priceElement.getAttribute('data-price'));
        if (!isNaN(price)) {
            const priceMinus15 = (price * 0.85).toFixed(2);
            const priceMinus15Element = priceElement.closest('tr').querySelector('[data-price-minus-15]');
            if (priceMinus15Element) {
                priceMinus15Element.textContent = priceMinus15;
            }
        }
    });
});