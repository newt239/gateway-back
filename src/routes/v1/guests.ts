import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db'
const router = express.Router();

router.get('/info/:guest_id', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const guest_id: string = req.params.guest_id;
    const sql: string = `SELECT * FROM gateway.guest WHERE guest_id=${guest_id}`;
    connection.query(sql, function (err: any, result: any) {
        if (err) {
            return res.json(err);
        } else {
            return res.json({
                status: "success",
                data: {
                    guest_id: result[0].guest_id,
                    guest_type: result[0].guest_type,
                    reservation_id: result[0].reservation_id,
                    place_id: result[0].place_id,
                    part: result[0].part,
                    available: result[0].available,
                    note: result[0].note,
                    regist_at: result[0].regist_at,
                }
            });
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
        sql += `INSERT INTO gateway.guest (guest_id, guest_type, reservation_id, userid, regist_at) VALUES (${guest_id} ${guest_type} ${reservation_id} ${res.locals.userid}, ${new Date()});`;
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
        sql += `UPDATE gateway.guest SET revoke_at=${new Date()} WHERE guest_id=${guest_id}`;
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