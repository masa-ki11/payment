document.getElementById("resetForm").addEventListener("submit", async (e) => {
    e.preventDefault();

    const password = document.getElementById("password").value;
    const password_confirm = document.getElementById("password_confirm").value;
    const token = new URLSearchParams(window.location.search).get("token");

    const response = await fetch("/reset", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ password, password_confirm, token }),
    });

    if (response.ok) {
        alert("パスワードがリセットされました。");
        window.location.href = "/login";
    } else {
        const errorData = await response.json();
        alert(`エラー: ${errorData.error}`);
    }
});
