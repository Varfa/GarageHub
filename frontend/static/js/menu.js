let activeMenu = null;
let activeButton = null;
let menuPlaceholder = null;

function closeActionMenu() {
    if (!activeMenu) {
        return;
    }

    activeMenu.classList.remove("is-open", "is-floating");
    activeMenu.removeAttribute("style");

    if (menuPlaceholder && menuPlaceholder.parentNode) {
        menuPlaceholder.parentNode.insertBefore(
            activeMenu,
            menuPlaceholder
        );

        menuPlaceholder.remove();
    }

    if (activeButton) {
        activeButton.setAttribute("aria-expanded", "false");
    }

    activeMenu = null;
    activeButton = null;
    menuPlaceholder = null;
}

function toggleMenu(button) {
    const menu = button.nextElementSibling;

    if (!menu || !menu.classList.contains("actions-menu")) {
        return;
    }

    if (activeMenu === menu) {
        closeActionMenu();
        return;
    }

    closeActionMenu();

    activeMenu = menu;
    activeButton = button;

    menuPlaceholder = document.createComment(
        "actions-menu-placeholder"
    );

    menu.parentNode.insertBefore(menuPlaceholder, menu);

    document.body.appendChild(menu);

    menu.classList.add("is-open", "is-floating");
    button.setAttribute("aria-expanded", "true");

    positionActionMenu(button, menu);
}

function positionActionMenu(button, menu) {
    const buttonRect = button.getBoundingClientRect();
    const viewportPadding = 12;
    const menuGap = 8;

    const menuWidth = menu.offsetWidth;
    const menuHeight = menu.offsetHeight;

    let left = buttonRect.right - menuWidth;
    let top = buttonRect.bottom + menuGap;

    if (left < viewportPadding) {
        left = viewportPadding;
    }

    if (left + menuWidth > window.innerWidth - viewportPadding) {
        left = window.innerWidth - menuWidth - viewportPadding;
    }

    const doesNotFitBelow =
        top + menuHeight >
        window.innerHeight - viewportPadding;

    if (doesNotFitBelow) {
        top = buttonRect.top - menuHeight - menuGap;
    }

    if (top < viewportPadding) {
        top = viewportPadding;
    }

    menu.style.left = `${Math.round(left)}px`;
    menu.style.top = `${Math.round(top)}px`;
}

document.addEventListener("click", (event) => {
    if (
        activeMenu &&
        !activeMenu.contains(event.target) &&
        event.target !== activeButton
    ) {
        closeActionMenu();
    }
});

document.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
        closeActionMenu();
    }
});

window.addEventListener("resize", closeActionMenu);
window.addEventListener("scroll", closeActionMenu, true);
