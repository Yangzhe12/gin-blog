$(document).ready(function(){
        $(".form-regist").submit(function(e){
          e.preventDefault();
          $(".error-msg").html("")
          phone = $("#phone").val()
          username = $("#username").val()
          password = $("#password").val()
          repwd = $("#repwd").val()
          email = $("#email").val()

          // console.log(phone)
          // 验证表单数据
          var phoneReg = /^(13[0-9]{1}|14[5|7|9]{1}|15[0-3|5-9]{1}|166|17[0-3|5-8]{1}|18[0-9]{1}|19[8-9]{1}){1}\d{8}$/;

          if (!phoneReg.test(phone)) {
            $("#phone-msg span").html("手机号格式错误，请检查！")
            return
          }

          var usernameReg = /^[a-zA-Z0-9_]{5,16}$/
          if(!usernameReg.test(username)) {
            if (username.length < 5) {
              $("#username-msg span").html("用户名长度为5-16位！")
            } else {
              $("#username-msg span").html("用户名由数字、大小写字母以及下划线成！")
            }
            return
          }

          if (password.length < 8) {
            $("#password-msg span").html("密码必须8位以上！")
          }

          if (repwd != password) {
            $("#password-msg span").html("两次密码不相同，请检查！")
          }

          var emailReg = /^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$/;
          if (!emailReg.test(email)) {
            $("#email-msg span").html("邮箱格式错误！")
          }

          var data = {
            phone: phone,
            username: username,
            password: password,
            repwd: repwd,
            email:email
          }

          var jsonData = JSON.stringify(data)
          var tokenData = $("#csrfToken").text()

          $.ajax({
            url:"/v1/regist",
            type: "post",
            data: jsonData,
            contentType:"application/json",
            headers:{
              "X-CSRF-TOKEN":tokenData,
            },
            success: function (resp) {
              if (resp.resno == "0") {
                alert(resp.msg)
                location.href = "/v1/login";
              } else {
                if (resp.resno == "1" | resp.resno == "2") {
                  alert(resp.msg)
                } else if (resp.resno == "3") {
                  $("#username-msg span").html(resp.msg)
                } else if (resp.resno == "4") {
                  $("#email-msg span").html(resp.msg)
                } else if (resp.resno == "5") {
                  alert(resp.msg)
                }else if (resp.resno == "10") {
                  alert(resp.msg)
                }
                
              }
            }
          })
        });
      });