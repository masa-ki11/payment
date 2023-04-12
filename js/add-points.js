function addPoints() {
    // ポイントの加算先となるユーザーIDを配列に格納する
    var userIds = [];
    var checkboxes = document.getElementsByName("user-id");
    for (var i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i].checked) {
            userIds.push(checkboxes[i].value);
        }
    }
    // console.log(userIds)
    // 加算するポイント数を取得する
    var points = document.getElementById("add-points").value;

    // ポイントを加算するためのAPIを呼び出す
    const requestBody = JSON.stringify({
        user_ids: userIds,
        points: parseInt(points, 10)
    });

    fetch("/add-points", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: requestBody
    }).then(function(response) {
        if (response.ok) {
            alert("ポイントを付与しました");
        } else {
            alert("errors");
        }
        location.reload(); // ページをリロードする
    });
}
document.addEventListener('DOMContentLoaded', function() {
    document.querySelector('.btn-add-points').addEventListener('click', addPoints);
});

function toggleSelectAll(source) {
    var checkboxes = document.getElementsByName("user-id");
    for (var i = 0; i < checkboxes.length; i++) {
        checkboxes[i].checked = source.checked;
    }
}
document.addEventListener('DOMContentLoaded', function() {
    var selectAllCheckbox = document.querySelector('input[name="select-all"]');
        selectAllCheckbox.addEventListener('click', function() {
            toggleSelectAll(selectAllCheckbox);
        });
});


document.addEventListener("DOMContentLoaded", function() {
    const pointElements = document.getElementsByClassName("point");
    for (let i = 0; i < pointElements.length; i++) {
        const pointValue = parseInt(pointElements[i].innerText, 10);
        pointElements[i].innerText = pointValue.toLocaleString();
    }
});

