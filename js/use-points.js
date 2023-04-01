$(document).on("submit", "#operation-form", function (event) {
    event.preventDefault(); // フォームのデフォルトの送信をキャンセルする

    const userId = $("#user-id").data("user-id"); // ユーザーIDを取得する
    const points = $("#point-input").val(); // ポイント数を取得する
    const action = $("#action-select").val(); // アクションを取得する
    const details = $("#details-input").val(); // 詳細を取得する

    // APIに送信するデータを作成する
    const requestData = {
        user_ids: [userId.toString()],
        points: parseInt(points),
        details: details
    };

    let apiUrl = '';
    if (action === '付与') {
        apiUrl = '/add-points';
    } else if (action === '使用') {
        apiUrl = '/use-points';
    }
console.log(requestData);
    // ポイントを操作するAPIにデータを送信する
    fetch(apiUrl, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    })
    .then(function(response) {
        if (response.ok) {
            alert("ポイントを操作しました。");
        } else {
            alert("ポイント操作に失敗しました。");
        }
        location.reload(); // ページをリロードする
    })
});
