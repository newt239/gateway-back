import { Connection } from 'mysql2';
const mysql = require('mysql2');
require('dotenv').config();

export const connectDb = (userId: string, password: string) => {
    const connection: Connection = mysql.createConnection({
        host: process.env.MYSQL_DATABASE_HOST,
        user: userId,
        password: password,
        database: "gateway",
        multipleStatements: true
    });
    return connection;
};