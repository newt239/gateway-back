const mysql = require('mysql2');
require('dotenv').config();

export const connectDb = (userid: string, password: string) => {
    const connection = mysql.createConnection({
        host: process.env.MYSQL_DATABASE_HOST,
        user: userid,
        password: password,
        database: "gateway"
    });
    return connection;
}