function confirmDelete(id) {
    if (confirm("削除しますか？")) {
        var form = document.createElement("form");
        form.setAttribute("method", "post");
        form.setAttribute("action", "/delete-user");

        var inputId = document.createElement("input");
        inputId.setAttribute("type", "hidden");
        inputId.setAttribute("name", "user-id");
        inputId.setAttribute("value", id);
        form.appendChild(inputId);

        document.body.appendChild(form);
        fetch('/delete-user', {
            method: 'POST',
            body: new FormData(form),
            headers: {
                'Accept': 'application/json',
            }
        })
        .then(response => response.json())
        .then(data => {
            if (data.status === "User deleted") {
                alert(data.message);
                location.reload();
            } else {
                alert("削除に失敗しました: " + data.error);
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert("削除に失敗しました");
        });
    }
}

document.addEventListener('DOMContentLoaded', function() {
    var deleteButtons = document.querySelectorAll('.btn-delete-user');

    deleteButtons.forEach(function(button) {
        button.addEventListener('click', function(event) {
            var userId = event.target.getAttribute('data-user-id');
            console.log(userId);
            confirmDelete(userId);
        });
    });
});

