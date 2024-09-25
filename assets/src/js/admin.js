document.querySelectorAll("a[href$='/delete']").forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();

        if (confirm('Are you sure you want to delete this item?')) {
            window.location.href = this.getAttribute('href');
        }
    });
});