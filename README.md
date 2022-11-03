GoLineMongo
===
以Golang為基底結合MongoDB及Line Go SDK做成的小專案

說明
---
* receive
	* 此為接收Line方傳送給我方的資訊，並將傳送訊息的使用者及傳送的文字寫入進DB
* sendMsgToLine
	* 傳送hello至指定的使用者 (userid可做更換)
* queryMsg
	* 從db列出所有使用者傳送的所有文字
* queryMsg/:name
	* 從db列出指定姓名使用者傳送的所有文字
