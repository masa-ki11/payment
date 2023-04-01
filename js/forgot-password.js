document.getElementById('forgotForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const email = document.getElementById('email').value;

    const response = await fetch('/forgot', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
    });

    if (response.ok) {
        alert('パスワードリセットメールが送信されました。');
    } else {
        const errorData = await response.json();
        alert(`エラー: ${errorData.error}`);
        if (errorData.error === 'ユーザーが見つかりませんでしたユーザー登録をしてください') {
            window.location.href = '/register';
        }
    }
});
