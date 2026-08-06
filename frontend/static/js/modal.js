document.addEventListener("DOMContentLoaded", () => {
    initializeDeleteConfirmation();
});

function initializeDeleteConfirmation() {
    const dialog = document.getElementById("confirm-delete-dialog");
    const message = document.getElementById("confirm-delete-message");
    const cancelButton = document.getElementById("confirm-delete-cancel");
    const submitButton = document.getElementById("confirm-delete-submit");

    if (
        !dialog ||
        !message ||
        !cancelButton ||
        !submitButton
    ) {
        return;
    }

    let pendingForm = null;

    document.addEventListener("submit", (event) => {
        const form = event.target;

        if (!(form instanceof HTMLFormElement)) {
            return;
        }

        const actionURL = new URL(
            form.action,
            window.location.origin
        );

        if (!actionURL.pathname.endsWith("/delete")) {
            return;
        }

        if (form.dataset.confirmed === "true") {
            delete form.dataset.confirmed;
            return;
        }

        event.preventDefault();

        pendingForm = form;
        message.textContent = getDeleteMessage(actionURL.pathname);

        dialog.showModal();
    });

    cancelButton.addEventListener("click", () => {
        pendingForm = null;
        dialog.close();
    });

    submitButton.addEventListener("click", () => {
        if (!pendingForm) {
            dialog.close();
            return;
        }

        const form = pendingForm;
        pendingForm = null;

        form.dataset.confirmed = "true";
        dialog.close();

        form.requestSubmit();
    });

    dialog.addEventListener("click", (event) => {
        if (event.target === dialog) {
            pendingForm = null;
            dialog.close();
        }
    });

    dialog.addEventListener("cancel", () => {
        pendingForm = null;
    });
}

function getDeleteMessage(pathname) {
    const messages = {
        "/clients/delete":
            "Клиент будет удалён из системы. Это действие нельзя отменить.",

        "/cars/delete":
            "Автомобиль будет удалён из системы. Это действие нельзя отменить.",
    };

    return (
        messages[pathname] ||
        "Запись будет удалена. Это действие нельзя отменить."
    );
}
