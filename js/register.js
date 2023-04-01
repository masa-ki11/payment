$("#register-form").submit(function(event) {
    event.preventDefault();

    const formData = {
        name: $("#name").val(),
        email: $("#email").val(),
        password: $("#password").val(),
        password_confirm: $("#password_confirm").val(),
    };

    $.ajax({
        type: "POST",
        url: "/register",
        data: JSON.stringify(formData),
        contentType: "application/json",
        success: function(response) {
        alert("User created successfully");
        window.location.href = "/login"; // ログインページへリダイレクト
    },
    error: function(response) {
    alert(response.responseJSON.error || response.responseJSON.message);
    },
    });
});