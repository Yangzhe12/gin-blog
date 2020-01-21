$(document).ready(function(){
  $(".form-signin").submit(function(e){
    e.preventDefault();
    $(".error-msg").html("")
    username = $("#inputUsername").val()
    password = $("#inputPassword").val()
    var usernameReg = /^[a-zA-Z0-9_]{5,16}$/
    if(!usernameReg.test(username)) {
      if (username.length < 5) {
        $("#error-div span").html("用户名长度为5-16位！")
      } else {
        $("#error-div span").html("用户名由数字、大小写字母以及下划线成！")
      }
      return
    }

    var data = {
      username: username,
      password: password,
    }

    var jsonData = JSON.stringify(data)
    var tokenData = $("#csrfToken").text()

    $.ajax({
      url:"/v1/login",
      type:"post",
      data:jsonData,
      headers:{
        "X-CSRF-TOKEN":tokenData,
      },
      contentType:"application/json",
      success: function (resp) {
        if (resp.resno == "0") {
          alert(resp.msg)
          location.href = "/v1";
        } else {
          if (resp.resno == "1") {
            alert(resp.msg)
          } else if (resp.resno == "2") {
            $("#error-div span").html(resp.msg)
          } else if (resp.resno == "3") {
            $("#error-div span").html(resp.msg)
          } else if (resp.resno == "5") {
            alert(resp.msg)
          } else if (resp.resno == "10") {
            alert(resp.msg)
          }
        }
      }
    });
  });
});