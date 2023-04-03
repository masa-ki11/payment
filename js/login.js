// js/login.js
console.log("log.js loaded");
document.getElementById("login-form").addEventListener("submit", async (event) => {
    event.preventDefault();

    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    const response = await fetch("/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
    });
    if (response.ok) {
        const data = await response.json();
        // JWT トークンを保存
        localStorage.setItem("jwt", data.jwt);
        // ログイン後の画面にリダイレクト
        window.location.href = "/";
    } else {
        // エラーメッセージを表示
        const error = await response.json();
        alert(error.error || "ログインに失敗しました");
    }
});

