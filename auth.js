document.addEventListener("DOMContentLoaded", function () {
    console.log("auth.js загружен");

    // === ОБРАБОТКА ФОРМЫ ЛОГИНА ===
    const loginForm = document.getElementById("auth-form");
    if (loginForm) {
        loginForm.addEventListener("submit", async function (event) {
            event.preventDefault();

            const email = document.getElementById("email").value;
            const password = document.getElementById("password").value;

            try {
                const response = await fetch("/api/login", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ email, password }),
                });

                const data = await response.json();
                console.log("Ответ сервера:", data);

                if (response.ok) {
                    alert("✅ Успешный вход!");
                    localStorage.setItem("token", data.token);
                    window.location.href = "index.html"; // Переход на главную
                } else {
                    alert("❌ Ошибка: " + data.message);
                }
            } catch (error) {
                console.error("Ошибка запроса:", error);
                alert("🚨 Ошибка соединения с сервером!");
            }
        });
    }

    // === ОБРАБОТКА ФОРМЫ РЕГИСТРАЦИИ ===
    const registerForm = document.getElementById("register-form");
    if (registerForm) {
        registerForm.addEventListener("submit", async function (event) {
            event.preventDefault();

            const name = document.getElementById("name").value;
            const email = document.getElementById("email").value;
            const password = document.getElementById("password").value;

            try {
                const response = await fetch("/api/register", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ name, email, password }),
                });

                const data = await response.json();
                console.log("Ответ сервера:", data);

                if (response.ok) {
                    alert("✅ Регистрация успешна! Теперь войдите в аккаунт.");
                    window.location.href = "login.html"; // Перенаправление на логин
                } else {
                    alert("❌ Ошибка: " + data.message);
                }
            } catch (error) {
                console.error("Ошибка запроса:", error);
                alert("🚨 Ошибка соединения с сервером!");
            }
        });
    }
});
