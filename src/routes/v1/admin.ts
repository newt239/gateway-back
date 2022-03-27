import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db';
const router = express.Router();

router.post('/create', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userId, res.locals.password);
    const userId: string = req.body.userId;
    let sql = `CREATE USER '${userId}'@'localhost' IDENTIFIED BY '${req.body.password}';`;
    switch (req.body.userType) {
        case "moderator":
            sql += `GRANT ALL ON gateway.* TO '${userId}'@'localhost' WITH MAX_QUERIES_PER_HOUR 1000;`;
            break;
        case "executive":
            sql += `GRANT INSERT, UPDATE, SELECT ON gateway.* TO '${userId}'@'localhost' WITH MAX_QUERIES_PER_HOUR 500;`;
            break;
        case "exhibit":
            sql += `GRANT INSERT, UPDATE, SELECT ON gateway.* TO '${userId}'@'localhost' WITH MAX_QUERIES_PER_HOUR 300;`;
            break;
        case "analysis":
            sql += `GRANT SELECT ON gateway.* TO '${userId}'@'localhost' WITH MAX_QUERIES_PER_HOUR 500;`;
            break;
        default:
            return ({ status: "error", message: "userType wes incorrect." });
    };
    sql += `FLUSH PRIVILEGES;`;
    sql += `INSERT INTO gateway.user (user_id, display_name, user_type, available) VALUES ('${userId}', '${req.body.displayName}', '${req.body.userType}', 1);`
    connection.query(sql, function (err: any, result: any) {
        if (err) {
            return res.json({
                status: "error",
                error: err
            });
        } else {
            return res.json({
                status: "success",
                data: result
            });
        };
    });
});

module.exports = router;