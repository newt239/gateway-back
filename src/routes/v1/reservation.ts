import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db';
import { QueryError } from 'mysql2';
const router = express.Router();

router.get('/:reservation_id', verifyToken, function (req: express.Request, res: express.Response) {
    const connection = connectDb(res.locals.userId, res.locals.password);
    const reservation_id: string = req.params.reservation_id;
    const sql: string = `SELECT * FROM gateway.reservation WHERE reservation_id='${reservation_id}'`;
    connection.query(sql, function (err: QueryError, result: any) {
        if (err) {
            return res.json(err);
        } else {
            console.log(result);
            if (result.length === 0) {
                return res.status(400).json({
                    message: `${reservation_id}という予約は存在しません。`
                });
            } else {
                return res.json(result[0]);
            };
        };
    });
    connection.end();
});

export default router;