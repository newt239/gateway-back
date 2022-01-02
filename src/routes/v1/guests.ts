import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db'
const router = express.Router();

router.get('/info/:guest_id', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const guest_id: string = req.params.guest_id;
    const sql: string = `SELECT * FROM gateway.guest WHERE guest_id='${guest_id}'`;
    connection.query(sql, function (err: any, result: any) {
        if (err) {
            return res.status(400).json(err);
        } else {
            if (result.length === 0) {
                return res.json({
                    status: "error",
                    message: `${guest_id}というゲストは存在しません。`
                });
            } else {
                return res.json({
                    status: "success",
                    data: result[0]
                });
            }
        };
    });
});

router.post('/regist', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    let sql: string = "";
    for (const eachRegist of req.body) {
        const guest_id: string = eachRegist.guest_id;
        const reservation_id: string | null = eachRegist.reservation_id;
        // TODO: reservation tableとapiを作り次第修正
        const guest_type: string = "student";
        const timestamp = new Date().toLocaleString('ja-JP').slice(0, 19).replace('T', ' ');
        sql += `INSERT INTO gateway.guest (guest_id, guest_type, reservation_id, userid, regist_at) VALUES ('${guest_id}' '${guest_type}' '${reservation_id}' '${res.locals.userid}', '${timestamp}');`;
    }
    console.log(sql);
    connection.query(sql, function (err: any, result: any) {
        if (err) {
            return res.json(err);
        } else {
            return res.json({
                status: "success"
            });
        };
    });
});

router.post('/revoke', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    let sql: string = "";
    for (const eachRevoke of req.body) {
        const guest_id: string = eachRevoke.guest_id;
        const timestamp = new Date().toLocaleString('ja-JP').slice(0, 19).replace('T', ' ');
        sql += `UPDATE gateway.guest SET revoke_at=${timestamp} WHERE guest_id='${guest_id}'`;
    }
    connection.query(sql, function (err: any, result: any) {
        if (err) {
            return res.json(err);
        } else {
            return res.json({
                status: "success"
            });
        };
    });
});

module.exports = router;