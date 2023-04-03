// カレントユーザーの ID を取得
const userId = $("#user-id").data("user-id");

fetch(`/get-points?user_id=${userId}`)
    .then(response => response.json())
    .then(data => {
        if (data.hasOwnProperty("point")) {
            document.getElementById("currentPoints").innerText = data.point.toLocaleString();
        } else {
            console.error("Failed to retrieve user's points.");
        }
    });
