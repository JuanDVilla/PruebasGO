<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
    <style>
        @import url(https://fonts.googleapis.com/css?family=Open+Sans:400,700);

        body {
        background: #456;
        font-family: 'Open Sans', sans-serif;
        }

        .login {
        width: 400px;
        margin: 16px auto;
        font-size: 16px;
        }

        /* Reset top and bottom margins from certain elements */
        .login-header,
        .login p {
        margin-top: 0;
        margin-bottom: 0;
        }

        /* The triangle form is achieved by a CSS hack */
        .login-triangle {
        width: 0;
        margin-right: auto;
        margin-left: auto;
        border: 12px solid transparent;
        border-bottom-color: #28d;
        }

        .login-header {
        background: #28d;
        padding: 20px;
        font-size: 1.4em;
        font-weight: normal;
        text-align: center;
        text-transform: uppercase;
        color: #fff;
        }

        .login-header2 {
        background: #28d;
        padding: 20px;
        font-size: 1em;
        font-weight: normal;
        text-align: center;
        text-transform: uppercase;
        color: #fff;
        }

        .login-container {
        background: #ebebeb;
        padding: 12px;
        }

        /* Every row inside .login-container is defined with p tags */
        .login p {
        padding: 12px;
        }

        .login input {
        box-sizing: border-box;
        display: block;
        width: 100%;
        border-width: 1px;
        border-style: solid;
        padding: 16px;
        outline: 0;
        font-family: inherit;
        font-size: 0.95em;
        }

        .login input[type="file"] {
        background: #fff;
        border-color: #bbb;
        color: #555;
        }

        /* Text fields' focus effect */
        .login input[type="file"]:focus {
        border-color: #888;
        }

        .login input[type="submit"] {
        background: rgb(107, 143, 173);
        border-color: transparent;
        color: rgb(0, 0, 0);
        cursor: pointer;
        }

        .login input[type="submit"]:hover {
        background: rgb(167, 191, 202);
        }

        /* Buttons' focus effect */
        .login input[type="submit"]:focus {
        border-color:rgb(255, 255, 255);
        }

    </style>
</head>

<body>
   
    <div class="login">
        <div class="login-triangle"></div>
        
        <h2 class="login-header">VALIDACIÓN DE ANTECEDENTES
      
            <form action="/POST" method="POST" enctype="multipart/form-data">
                <h4 class="login-header2">Documento Parte Frontal
                    <input type="file" name="Frontal">
                </h4>
                <hr>
                <h4 class="login-header2">Documento Parte Reverso
                    <input type="file" name="Reverso">
                </h4>
                <p><input type="submit" value="VALIDAR ANTECEDENTES"></p>
            </form>
        </h2>
      </div>
</body>

</html>