module.exports = {
    extends: ["@commitlint/config-conventional"],
    rules: {
        "subject-max-length": [2, "always", 72], // 标题最多 72 个字符
    },
};
