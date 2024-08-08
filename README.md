# vcs-sms-microservice

## 1. Overview hệ thống:
![overview](pic/overview2.png)
- Hệ thống được phân chia thành các servieces như:
  * Healthcheck-server:để nhận request từ các agents.
  * Healthcheck-worker:thực hiện xử lý các request từ **healthcheck-server**, cập nhật các server vừa gửi thông tin, các server hiện thời.
  * Mail: thực hiện xử lý request report từ server-management gửi mail.
  * Server-management: thực hiện CRUD, gửi report hàng ngày.
### 1.1. Healthcheck-server:
![healthcheck_server](pic/healthcheck_server.png)
### 1.2. Healthcheck-worker:
![healthcheck_worker](pic/healthcheck_worker.png)
### 1.3. Server-management:
![Server-management](pic/Server_management.png)
### 1.4. Mail:
![mail](pic/mail.png)
