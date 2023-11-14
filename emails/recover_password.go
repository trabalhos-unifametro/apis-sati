package emails

import (
	"apis-sati/models"
	"apis-sati/utils"
	"fmt"
	"strings"
)

func SendEmailCodeRandom(user models.User) error {
	code := strings.Split(user.CodeRecovery, "")

	body := fmt.Sprint(`<p style="font-size: 20px"><b>Olá</b> `, user.Name, `!</p>
	   <p style="font-size: 20px">Segue abaixo o código para recuperar sua senha:</p>
	   <table style="width: 500px; height: 50px; background-color: #1C6F77; border-radius: 5px">
		   <tr align="center" valign="middle" style="color: #FFFFFF; font-size: 20px; font-weight: 700">
			   <td>`, code[0], `</td>
			   <td>`, code[1], `</td>
			   <td>`, code[2], `</td>
			   <td>-</td>
			   <td>`, code[3], `</td>
			   <td>`, code[4], `</td>
			   <td>`, code[5], `</td>
			   <td>`, code[6], `</td>
		   </tr>
	   </table>
	   <p style="font-size: 20px">
		   <b>Importante:</b> O código só valerá por 5 minutos, após esse período, <br/>
		   será necessário solicitar um novo código.
	   </p>
	   <p style="font-size: 15px">
		   Caso você não tenha realizado essa solicitação, por favor entre em contato conosco.<br/>
		   S.A.T.I nunca solicita que você informe sua senha por e-mail.<br/>
		   Não responda esta mensagem.
	   </p>
	   <p style="font-weight: bold; font-size: 18px">
		   Não compartilhe esse código com ninguém.<br/>
		   Nunca compartilhe sua senha com outra pessoa.
    </p>`)

	template := MountLayoutTemplateEmail("Código para recuperação de senha", body)
	err := MountEmail(template, "[S.A.T.I] - Código de verificação.", utils.CheckToSend(user.Email))

	return err
}

func SuccessfulRecoverPassword(user models.User) error {
	body := fmt.Sprint(`<p style="font-size: 25px"><b>Olá</b> `, user.Name, `!</p>
	   <p style="font-size: 25px">Sua senha foi recuperada com sucesso!</p>
	   <p style="font-size: 25px; line-height: normal">
		   Caso você não tenha realizado essa ação, por favor entre em contato conosco.<br/>
		   S.A.T.I nunca solicita que você informe sua senha por e-mail.<br/>
		   Não responda esta mensagem.
	   </p>
	`)

	template := MountLayoutTemplateEmail("Senha recuperada com sucesso.", body)
	err := MountEmail(template, "[S.A.T.I] - Recuperação de senha.", utils.CheckToSend(user.Email))

	return err
}

func MountLayoutTemplateEmail(title, body string) string {
	return fmt.Sprint(`<!-- Inliner Build Version 4380b7741bb759d6cb997545f3add21ad48f010b -->
	<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
	<html lang="pt" xmlns="http://www.w3.org/1999/xhtml" xmlns="http://www.w3.org/1999/xhtml">
	<head>
		<title></title>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
		<meta name="viewport" content="width=device-width"/>
		<link rel="preconnect" href="https://fonts.googleapis.com">
		<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
		<link href="https://fonts.googleapis.com/css2?family=Assistant:wght@400;700&family=Poppins:wght@400;600;700&display=swap" rel="stylesheet">
	</head>
	<body style="width: 100% !important;
				height: 100vh;
				-webkit-text-size-adjust: 100%;
				-webkit-font-smoothing: antialiased;
				-ms-text-size-adjust: 100%;
				color: #222222;
				font-family: 'Helvetica Neue', 'Arial', sans-serif;
				font-weight: normal; text-align: left;
				line-height: 19px;
				font-size: 14px;
				margin: 0;
				padding: 0;
				background-color: #F3F3F3;
	">
	<table style="width: 100%; height: 100%; background-color: transparent;" border="0" cellspacing="0">
		<tr valign="middle" style="height: 130px; background-color: #1C6F77; box-shadow: -1px 9px 13px 5px rgba(0, 0, 0, 0.25);">
			<th style="width: 230px; padding-left: 55px">
				<img src="https://portal-sati.s3.sa-east-1.amazonaws.com/emails/logo.png" width="330" alt="Hospital Medical Care" style="display: inline-block" />
			</th>
			<th style="width: 1px; padding: 0 30px">
				<div style="height: 90px; width: 1px; background-color: rgba(255, 255, 255, 0.5); display: inline-block"></div>
			</th>
			<th>
				<span style="font-size: 25px; color: #FFFFFF">`, title, `</span>
			</th>
		</tr>
		<tr>
		   <td colspan="3" align="center" valign="middle" style="padding: 30px 0;">
			   `, body, `
		   </td>
		</tr>
		<tr style="height: 170px;" align="center" valign="middle">
			<td colspan="3">
				<div style="width: 100%; border-top: 1px solid rgba(0,0,0,0.3); padding: 30px 0">
					<span>Caso tenha dúvidas, fale conosco através dos canais:</span><br/>
					<table style="margin: 20px 0;">
						<tr>
							<td>
								<img style="height: 20px; margin-bottom: -3px; margin-right: 5px; display: inline-block" src="https://portal-sati.s3.sa-east-1.amazonaws.com/emails/icon_telephone.png" alt="Telefone"/>
								<b style="font-size: 20px">(00) 0000-0000</b>
							</td>
							<td style="padding: 0 25px">
								<img style="height: 20px; margin-bottom: -3px; margin-right: 5px; display: inline-block" src="https://portal-sati.s3.sa-east-1.amazonaws.com/emails/icon_whatsapp.png" alt="Telefone"/>
								<b style="font-size: 20px">(00) 0 0000-0000</b>
							</td>
							<td>
								<img style="height: 20px; margin-bottom: -3px; margin-right: 5px; display: inline-block" src="https://portal-sati.s3.sa-east-1.amazonaws.com/emails/icon_email.png" alt="Telefone"/>
								<b style="font-size: 20px">suporte@sati.org.br</b>
							</td>
						</tr>
					</table>
					<span>Copyright © 2023 Hospital Medical Care</span>
				</div>
			</td>
		</tr>
	</table>
	</body>
	</html>`)
}
