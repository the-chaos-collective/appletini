import dotenv from 'dotenv';
dotenv.config();

// .graphqlrc.ts or graphql.config.ts
export default {
  projects: {
    appletini: {
      extensions: {
        endpoints: {
          default: {
            url: "https://api.github.com/graphql",
            headers : 
            {
              Authorization: `Bearer ${process.env.GITHUB_ACCESS_TOKEN}`
            }
          },
        },
      },
    },
  },
};