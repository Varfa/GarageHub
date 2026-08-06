document.addEventListener("DOMContentLoaded", () => {
    document.documentElement.classList.add("js-enabled");

    initializeActiveNavigation();
    initializePageTitle();
    initializeDashboardHeader();
    initializeMobileNavigation();
    initializeClickableRows();
});

function initializeActiveNavigation() {
    const currentPath = window.location.pathname;
    const navigationLinks = document.querySelectorAll(".sidebar a");

    navigationLinks.forEach((link) => {
        const linkURL = new URL(link.href, window.location.origin);
        const linkPath = linkURL.pathname;

        const isDashboard =
            linkPath === "/dashboard" &&
            currentPath === "/dashboard";

        const isModulePage =
            linkPath !== "/dashboard" &&
            currentPath.startsWith(linkPath);

        if (isDashboard || isModulePage) {
            link.classList.add("is-active");
            link.setAttribute("aria-current", "page");
        }
    });
}

function initializePageTitle() {
    const pageHeading = document.querySelector(".content h1");

    if (!pageHeading) {
        return;
    }

    const headingText = pageHeading.textContent.trim();

    if (!headingText) {
        return;
    }

    document.title = `${headingText} — GarageHub`;
}

function initializeDashboardHeader() {
    const greetingElement = document.getElementById("dashboard-greeting");
    const weekdayElement = document.getElementById("dashboard-weekday");
    const dateElement = document.getElementById("dashboard-date");

    if (!greetingElement || !weekdayElement || !dateElement) {
        return;
    }

    const now = new Date();
    const hour = now.getHours();
    const locale = getCurrentLocale();

    greetingElement.textContent = `${getGreeting(hour, locale)}, Roman`;

    weekdayElement.textContent = capitalizeFirstLetter(
        new Intl.DateTimeFormat(locale, {
            weekday: "long",
        }).format(now)
    );

    dateElement.textContent = new Intl.DateTimeFormat(locale, {
        day: "numeric",
        month: "long",
        year: "numeric",
    }).format(now);
}

function initializeMobileNavigation() {
    const header = document.querySelector(".app-header");
    const sidebar = document.querySelector(".sidebar");

    if (!header || !sidebar) {
        return;
    }

    const menuButton = document.createElement("button");
    menuButton.className = "mobile-menu-button";
    menuButton.type = "button";
    menuButton.setAttribute("aria-label", "Открыть меню");
    menuButton.setAttribute("aria-expanded", "false");

    menuButton.innerHTML = `
        <svg viewBox="0 0 24 24" aria-hidden="true">
            <path d="M4 6h16v2H4V6Zm0 5h16v2H4v-2Zm0 5h16v2H4v-2Z"></path>
        </svg>
    `;

    header.insertBefore(menuButton, header.firstChild);

    const overlay = document.createElement("div");
    overlay.className = "mobile-menu-overlay";
    overlay.setAttribute("aria-hidden", "true");

    document.body.appendChild(overlay);

    const openMenu = () => {
        sidebar.classList.add("is-open");
        overlay.classList.add("is-open");
        document.body.classList.add("mobile-menu-open");

        menuButton.setAttribute("aria-expanded", "true");
        menuButton.setAttribute("aria-label", "Закрыть меню");
    };

    const closeMenu = () => {
        sidebar.classList.remove("is-open");
        overlay.classList.remove("is-open");
        document.body.classList.remove("mobile-menu-open");

        menuButton.setAttribute("aria-expanded", "false");
        menuButton.setAttribute("aria-label", "Открыть меню");
    };

    menuButton.addEventListener("click", () => {
        if (sidebar.classList.contains("is-open")) {
            closeMenu();
            return;
        }

        openMenu();
    });

    overlay.addEventListener("click", closeMenu);

    sidebar.querySelectorAll("a").forEach((link) => {
        link.addEventListener("click", closeMenu);
    });

    document.addEventListener("keydown", (event) => {
        if (event.key === "Escape") {
            closeMenu();
        }
    });

    window.addEventListener("resize", () => {
        if (window.innerWidth > 720) {
            closeMenu();
        }
    });
}

function initializeClickableRows() {
    const rows = document.querySelectorAll("[data-href]");

    rows.forEach((row) => {
        const navigate = () => {
            const target = row.dataset.href;

            if (target) {
                window.location.href = target;
            }
        };

        row.addEventListener("click", (event) => {
            const interactiveElement = event.target.closest(
                "a, button, input, select, textarea, form"
            );

            if (interactiveElement) {
                return;
            }

            navigate();
        });

        row.addEventListener("keydown", (event) => {
            if (event.key !== "Enter" && event.key !== " ") {
                return;
            }

            if (event.target !== row) {
                return;
            }

            event.preventDefault();
            navigate();
        });
    });
}

function getCurrentLocale() {
    const languageCode = document.documentElement.lang || "en";

    const locales = {
        en: "en-GB",
        lt: "lt-LT",
        uk: "uk-UA",
        ru: "ru-RU",
    };

    return locales[languageCode] || "en-GB";
}

function getGreeting(hour, locale) {
    const language = locale.split("-")[0];

    const greetings = {
        en: {
            morning: "Good morning",
            afternoon: "Good afternoon",
            evening: "Good evening",
            night: "Good night",
        },
        lt: {
            morning: "Labas rytas",
            afternoon: "Laba diena",
            evening: "Labas vakaras",
            night: "Labos nakties",
        },
        uk: {
            morning: "Доброго ранку",
            afternoon: "Добрий день",
            evening: "Добрий вечір",
            night: "Доброї ночі",
        },
        ru: {
            morning: "Доброе утро",
            afternoon: "Добрый день",
            evening: "Добрый вечер",
            night: "Доброй ночи",
        },
    };

    const selectedGreetings = greetings[language] || greetings.en;

    if (hour >= 5 && hour < 12) {
        return selectedGreetings.morning;
    }

    if (hour >= 12 && hour < 18) {
        return selectedGreetings.afternoon;
    }

    if (hour >= 18 && hour < 23) {
        return selectedGreetings.evening;
    }

    return selectedGreetings.night;
}

function capitalizeFirstLetter(value) {
    if (!value) {
        return value;
    }

    return value.charAt(0).toUpperCase() + value.slice(1);
}
