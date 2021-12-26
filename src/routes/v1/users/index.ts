import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db'
var router = express.Router();

router.get('/me', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const sql = `SELECT * FROM gateway.user WHERE userid='${res.locals.userid}'`;
    connection.query(sql, function (err: any, result: any) {
        if (err) {
            return res.json(err);
        } else {
            return res.json({
                userid: result[0].userid,
                display_name: result[0].display_name,
                user_type: result[0].user_type,
                role: result[0].role,
                available: result[0].available,
                note: result[0].note,
            })
        };
    });
});

module.exports = router;