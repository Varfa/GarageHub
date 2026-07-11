function toggleMenu(button) {
    const menu = button.nextElementSibling;

    document.querySelectorAll(".actions-menu").forEach(item => {
        if (item !== menu) {
            item.style.display = "none";
        }
    });

    menu.style.display = menu.style.display === "block" ? "none" : "block";
}

document.addEventListener("click", function (event) {
    if (!event.target.closest(".actions-cell")) {
        document.querySelectorAll(".actions-menu").forEach(menu => {
            menu.style.display = "none";
        });
    }
});
function toggleCarForm() {
    const form = document.getElementById("car-create-form");

    form.classList.toggle("is-open");
}
