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
        // メール送信用
        // alert('パスワードリセットメールが送信されました。');

        // メール送信不可の代替案
        const data = await response.json();
        const resetUrl = data.resetUrl;
        window.location.href = resetUrl;
    } else {
        const errorData = await response.json();
        alert(`エラー: ${errorData.error}`);
        if (errorData.error === 'ユーザーが見つかりませんでしたユーザー登録をしてください') {
            window.location.href = '/register';
        }
    }
});
