# IdentityReconciliation
Bitespeed Backend Task: Identity Reconciliation

The Repo has the code of the task assigned for hiring to the position of Bitespeed Backend developer.

Link to Kiran H resume 
```
https://drive.google.com/file/d/1QM_CfB-W6y2j4KnGmEfD6B2_L0ZXC7F2/view?usp=drive_link
```

Tech stack used

Backend : Golang
Database : sqlite3 

steps to test out the /identify endpoint

1. open terminal and Clone the git repo 
```
https://github.com/KiranH98/IdentityReconciliation.git
```

2. go inside the project directory and build the docker 
```
docker build -t identity-reconciliation-app . 
```

3. now run the built docker 
```
 docker run -p 8080:8080 identity-reconciliation-app
```

4. The application has been added with swagger UI to test out API apart from postman , Link to swagger UI
```
http://localhost:8080/swagger/index.html
```

