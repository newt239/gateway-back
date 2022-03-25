import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db';
import { QueryError } from 'mysql2';
const router = express.Router();
const jwt = require('jsonwebtoken');

router.post('/login', function (req: express.Request, res: express.Response) {
    const userid: string = req.body.userid;
    const password: string = req.body.password;
    const connection = connectDb(userid, password);
    connection.connect(function (err) {
        if (err) {
            res.json({
                status: "error",
                message: "username or password were incorrect."
            });
            return;
        };
    });
    const token = jwt.sign({ userid: userid, password: password }, process.env.SIGNATURE);
    const sql = `SELECT * FROM gateway.user WHERE userid='${userid}'`;
    connection.query(sql, function (err: QueryError, result: any) {
        if (err) {
            res.json(err);
        } else {
            return res.json({
                status: "success",
                token: token,
                profile: {
                    userid: result[0].userid,
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
    const connection = connectDb(res.locals.userid, res.locals.password);
    const sql = `SELECT * FROM gateway.user WHERE userid='${res.locals.userid}'`;
    connection.query(sql, function (err: QueryError, result: any) {
        if (err) {
            console.log(res)
            return res.json(err);
        } else {
            return res.json({
                status: "success",
                profile: {
                    userid: result[0].userid,
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