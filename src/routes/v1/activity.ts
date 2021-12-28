import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db'
const router = express.Router();

router.post('/enter', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const guest_id: string = req.body.guest_id;
    const place_id: string = req.body.place_id;
    const timestamp = new Date();
    const activity_id: string = timestamp.getTime().toString(16) + Math.floor(Math.random() * 10).toString(16);
    const sql: string = `INSERT * FROM gateway.activity (activity_id, guest_id, place_id, userid, activity_type, timestamp, available) VALUES (${activity_id}, ${guest_id}, ${place_id}, ${res.locals.userid}, enter, ${timestamp}, 1)`;
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

router.post('/exit', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const guest_id: string = req.body.guest_id;
    const place_id: string = req.body.place_id;
    const timestamp = new Date();
    const activity_id: string = timestamp.getTime().toString(16) + Math.floor(Math.random() * 10).toString(16);
    const sql: string = `INSERT * FROM gateway.activity (activity_id, guest_id, place_id, userid, activity_type, timestamp, available) VALUES (${activity_id}, ${guest_id}, ${place_id}, ${res.locals.userid}, exit, ${timestamp}, 1)`;
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

router.post('/pass', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const guest_id: string = req.body.guest_id;
    const place_id: string = req.body.place_id;
    const timestamp = new Date();
    const activity_id: string = timestamp.getTime().toString(16) + Math.floor(Math.random() * 10).toString(16);
    const sql: string = `INSERT * FROM gateway.activity (activity_id, guest_id, place_id, userid, activity_type, timestamp, available) VALUES (${activity_id}, ${guest_id}, ${place_id}, ${res.locals.userid}, pass, ${timestamp}, 1)`;
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