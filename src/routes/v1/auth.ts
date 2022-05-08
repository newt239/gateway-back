import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db';
import { QueryError } from 'mysql2';
const router = express.Router();
const jwt = require('jsonwebtoken');

type userTypeProp =
    | "admin"
    | "moderator"
    | "executive"
    | "exhibit"
    | "analysis"
    | "temporary";

interface sqlAuthLoginResultProp {
    user_id: string;
    display_name: string;
    user_type: userTypeProp;
    role: string;
    available: number;
    note: string;
};

router.post('/login', function (req: express.Request, res: express.Response) {
    const userId: string = req.body.userId;
    const password: string = req.body.password;
    const connection = connectDb(userId, password);
    connection.connect((err) => {
        if (err) {
            res.status(400).json({
                message: "username or password were incorrect."
            });
            return;
        };
    });
    const token = jwt.sign({ userId: userId, password: password }, process.env.SIGNATURE);
    const sql = `SELECT * FROM gateway.user WHERE user_id='${userId}'`;
    connection.query(sql, (err: QueryError, result: sqlAuthLoginResultProp[]) => {
        if (err || result.length == 0) {
            res.json(err);
        } else {
            return res.json({
                token: token,
                profile: {
                    user_id: result[0].user_id,
                    display_name: result[0].display_name,
                    user_type: result[0].user_type,
                    role: result[0].role,
                    available: result[0].available === 0 ? true : false,
                    note: result[0].note
                }
            });
        };
    });
    connection.end();
});

router.get('/me', verifyToken, function (_req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userId, res.locals.password);
    const sql = `SELECT * FROM gateway.user WHERE user_id='${res.locals.userId}'`;
    connection.query(sql, function (err: QueryError, result: sqlAuthLoginResultProp[]) {
        if (err) {
            return res.status(400).json(err);
        } else {
            return res.json({
                status: "success",
                profile: {
                    user_id: result[0].user_id,
                    display_name: result[0].display_name,
                    user_type: result[0].user_type,
                    role: result[0].role,
                    available: result[0].available === 0 ? true : false,
                    note: result[0].note
                }
            });
        };
    });
    connection.end();
});

export default router;