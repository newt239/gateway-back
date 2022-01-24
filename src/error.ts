import { QueryError } from 'mysql2';

const triageError = (err: QueryError) => {
    switch (err.code) {
        case "ER_CON_COUNT_ERROR":
            return {
                status: "error",
                message: "to many connections"
            }
    };
};

export default triageError;