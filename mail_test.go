// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"testing"
)

//func TestMail(t *testing.T) {
//	config := `{"username":"359851485@qq.com","password":"Aizhu1218","host":"smtp.qq.com","port":587}`
//	mail := NewEMail(config)
//	if mail.Username != "359851485@qq.com" {
//		t.Fatal("email parse get username error")
//	}
//	if mail.Password != "Aizhu1218" {
//		t.Fatal("email parse get password error")
//	}
//	if mail.Host != "smtp.qq.com" {
//		t.Fatal("email parse get host error")
//	}
//	if mail.Port != 587 {
//		t.Fatal("email parse get port error")
//	}
//	mail.To = []string{"wwm86@126.com"}
//	mail.From = "Bigwish<359851485@qq.com>"
//	mail.Subject = "hi, just from beego! normal"
//	mail.Text = "Text Body is, of course, supported!"
//	//	mail.HTML = "<h1>Fancy Html is supported, too!</h1>"
//	mail.AttachFile("./safemap.go")
//	mail.Send()
//}

func TestSendUsingTLS(t *testing.T) {
	config := `{"username":"support@oceanelec.cn","password":"Ocean123","host":"smtp.exmail.qq.com","port":465}`
	mail := NewEMail(config)

	mail.To = []string{"wwm86@126.com"}

	mail.Subject = "hi, just from test!"
	mail.Text = "Text Body is, of course, supported! ssl"
	//	mail.HTML = "<h1>Fancy Html is supported, too!</h1>"
	mail.AttachFile("./safemap.go")
	mail.SendUsingTLS()
}
