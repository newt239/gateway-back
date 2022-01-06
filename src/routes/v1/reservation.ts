import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db';
const router = express.Router();

router.get('/:reservation_id', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userid, res.locals.password);
    const reservation_id: string = req.params.reservation_id;
    const sql: string = `SELECT * FROM gateway.reservation WHERE reservation_id='${reservation_id}'`;
    connection.query(sql, function (err: any, result: any) {
        if (err) {
            return res.json(err);
        } else {
            if (result.length === 0) {
                return res.json({
                    status: "error",
                    message: `${reservation_id}という予約は存在しません。`
                });
            } else {
                return res.json({
                    status: "success",
                    data: {
                        reservation_id: result[0].reservation_id,
                        guest_type: result[0].guest_type,
                        part: result[0].part,
                        available: result[0].available,
                        count: result[0].count,
                        note: result[0].note
                    }
                });
            };
        };
    });
});

module.exports = router;