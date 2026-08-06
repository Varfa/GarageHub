function toggleCarForm() {
    const form = document.getElementById("car-create-form");

    if (!form) {
        return;
    }

    form.classList.toggle("is-open");

    if (form.classList.contains("is-open")) {
        const firstInput = form.querySelector(
            'input:not([type="hidden"]), select, textarea'
        );

        if (firstInput) {
            firstInput.focus();
        }
    }
}
