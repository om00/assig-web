# Proejct - This simple web page project where admin can block and unlock the user
Admin can serach the user with name ,email , phoneNumber 
Mutltiple phone Number can be  used filter comma serpated

## local Installation
1 - golang should be installed in local sytem.  
2 - postgreSql database server should be installed.  
3 - create database which you want use for project or use any existing.  
4 - extract the zip-file in your system where want run the project.  
5 - go inside directry and update db configuration in .env file according to your system.  
6 - open the command line .
    run the command to go mod tidy -# (optional) downloads missing deps
    run the command to with migration - go run main.go -m up - this will create table in db.  
    7 - Now to create produre and function db that are list in  assig-web/prodfunc/ - folder 
    -create all the producre and functions manully in db.  
    8   Now to run the project run below command 
    go run main.go -s seedUsers  -#this will populate same sample data   into  user table and will start the project
    on 8085 port 
    - -s seedUsers use only for first time for sample data other if need to up project agian then you use simple 
    go run main.go.  
    9  To accssee webage open in you browser http://localhost:8085/  