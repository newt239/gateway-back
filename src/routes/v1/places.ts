import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db'
const router = express.Router();

router.get('/places/heat', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const guest_id: string = req.params.guest_id;
    const sql: string = `SELECT * FROM gateway.guest WHERE guest_id=${guest_id}`;
    connection.query(sql, function (err: any, result: any) {
        if (err) {
            return res.json(err);
        } else {
            return res.json({
                status: "success",
                data: {}
            });
        };
    });
});

router.get('/places/info/:place_id', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const place_id: string = req.params.place_id;
    const sql: string = `SELECT * FROM gateway.place WHERE place_id=${place_id}`;
    connection.query(sql, function (err: any, result: any) {
        if (err) {
            return res.json(err);
        } else {
            return res.json({
                status: "success",
                data: {
                    place_id: result[0].place_id,
                }
            });
        };
    });
});

module.exports = router;