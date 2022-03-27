import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db';
import { QueryError } from 'mysql2';
const router = express.Router();
const jwt = require('jsonwebtoken');

router.post('/login', function (req: express.Request, res: express.Response) {
    const userId: string = req.body.userId;
    const password: string = req.body.password;
    const connection = connectDb(userId, password);
    connection.connect(function (err) {
        if (err) {
            res.json({
                status: "error",
                message: "username or password were incorrect."
            });
            return;
        };
    });
    const token = jwt.sign({ userId: userId, password: password }, process.env.SIGNATURE);
    const sql = `SELECT * FROM gateway.user WHERE user_id='${userId}'`;
    connection.query(sql, function (err: QueryError, result: any) {
        if (err) {
            res.json(err);
        } else {
            return res.json({
                status: "success",
                token: token,
                profile: {
                    userId: result[0].user_id,
                    display_name: result[0].display_name,
                    user_type: result[0].user_type,
                    role: result[0].role,
                    available: result[0].available,
                    note: result[0].note
                }
            });
        };
    });
    connection.end();
});

router.get('/me', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userId, res.locals.password);
    const sql = `SELECT * FROM gateway.user WHERE user_id='${res.locals.userId}'`;
    connection.query(sql, function (err: QueryError, result: any) {
        if (err) {
            console.log(res)
            return res.json(err);
        } else {
            return res.json({
                status: "success",
                profile: {
                    userId: result[0].user_id,
                    display_name: result[0].display_name,
                    user_type: result[0].user_type,
                    role: result[0].role,
                    available: result[0].available,
                    note: result[0].note
                }
            });
        };
    });
    connection.end();
});

module.exports = router;