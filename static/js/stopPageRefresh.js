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
