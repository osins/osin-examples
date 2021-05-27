/*
 *
 * login-register modal
 * Autor: Creative Tim
 * Web-autor: creative.tim
 * Web script: http://creative-tim.com
 * 
 */
function showRegisterForm() {
    $('.loginBox').fadeOut('fast', function () {
        $('.registerBox').fadeIn('fast');
        $('.login-footer').fadeOut('fast', function () {
            $('.register-footer').fadeIn('fast');
        });
        $('.modal-title').html('登录');
    });

    $('#response_type').val('register');
    $('.error').removeClass('alert alert-danger').html('');

}
function showLoginForm() {
    $('#loginModal .registerBox').fadeOut('fast', function () {
        $('.loginBox').fadeIn('fast');
        $('.register-footer').fadeOut('fast', function () {
            $('.login-footer').fadeIn('fast');
        });

        $('.modal-title').html('登录');
    });
    $('.error').removeClass('alert alert-danger').html('');
}

function openLoginModal() {
    showLoginForm();
    setTimeout(function () {
        $('#loginModal').modal('show');
    }, 230);

}
function openRegisterModal() {
    showRegisterForm();
    setTimeout(function () {
        $('#loginModal').modal('show');
    }, 230);

}

function loginSubmit() {
    $('#response_type').val("login")
    var data = {
        username: $('#username').val(),
        password: $('#password').val(),
        client_id: $('#client_id').val(),
        client_secret: $('#client_secret').val(),
        response_type: $('#response_type').val(),
        redirect_uri: $('#redirect_uri').val(),
        state: $('#state').val(),
    };

    console.log(data);

    $.post("/oauth/authorize", data)
        .done(function (data) {
            if (data) {
                window.location.href = data.redirect_uri;
            }
        }).fail(function (e) {
            console.log(e.status == 500);
            shakeModal("用户名或密码错误.")
        });
}

function registerSubmit() {
    $('#response_type').val("register")
    var data = {
        username: $('#reg_username').val(),
        password: $('#reg_password').val(),
        email: $('#email').val(),
        mobile: $('#mobile').val(),
        client_id: $('#client_id').val(),
        client_secret: $('#client_secret').val(),
        response_type: $('#response_type').val(),
        redirect_uri: $('#redirect_uri').val(),
        state: $('#state').val(),
    };

    console.log(data);

    $.post("/oauth/authorize", data)
        .done(function (data) {
            if (data) {
                window.location.href = data.redirect_uri;
            }
        }).fail(function (e) {
            console.log(e.status == 500);
            shakeModal("用户名或密码错误.")
        });
}

function shakeModal(msg) {
    $('#loginModal .modal-dialog').addClass('shake');
    $('.error').addClass('alert alert-danger').html(msg);
    $('input[type="password"]').val('');
    setTimeout(function () {
        $('#loginModal .modal-dialog').removeClass('shake');
    }, 1000);
}

