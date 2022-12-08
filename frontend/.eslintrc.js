module.exports = {
    "env": {
        "browser": true,
        "es2021": true
    },
    "extends": [
        "eslint:recommended",
        "plugin:react/recommended",
        "plugin:@typescript-eslint/recommended",
        "next/core-web-vitals"
    ],
    "parser": "@typescript-eslint/parser",
    "parserOptions": {
        "ecmaFeatures": {
            "jsx": true
        },
        "ecmaVersion": 12,
        "sourceType": "module"
    },
    "plugins": [
        "react",
        "@typescript-eslint"
    ],
    "rules": {
        // suppress errors for missing 'import React' in files
       "react/react-in-jsx-scope": "off",
        // allow jsx syntax in js + tsx files
       "react/jsx-filename-extension": [1, { "extensions": [".js", ".jsx", ".tsx"] }],
       "@typescript-eslint/ban-ts-comment": "off"
    }
};
