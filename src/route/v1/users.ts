import express from 'express';

const router = express.Router();

// GETリクエスト
router.get('/', (req: express.Request, res: express.Response) => {
    try {
        res.status(200).json({ userId: "U001", userName: "Yamada Taro" });
    } catch (error) {
        res.status(400).json({ message: error });
    }
});

// POSTリクエスト
router.post('/', (req: express.Request, res: express.Response) => {
    try {
        res.status(200).json({ message: "登録しました" });
    } catch (error) {
        res.status(400).json({ message: error });
    }
});

export default router;