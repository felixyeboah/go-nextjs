This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

# Frontend Application

A modern frontend application built with Next.js, featuring:

- Modern React with Next.js App Router
- Type-safe with TypeScript
- Tailwind CSS for styling with shadcn/ui components
- Form validation with Zod and React Hook Form
- Data fetching with React Query
- Authentication with JWT
- Responsive design

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Project Structure

```
frontend/
├── app/                  # Next.js app directory
│   ├── api/              # API routes
│   ├── auth/             # Authentication pages
│   └── users/            # User pages
├── public/               # Static assets
├── src/                  # Source code
│   ├── components/       # React components
│   │   ├── auth/         # Authentication components
│   │   ├── layout/       # Layout components
│   │   └── ui/           # UI components
│   ├── lib/              # Utility functions
│   │   ├── api/          # API client
│   │   └── utils/        # Helper functions
│   └── types/            # TypeScript type definitions
├── .env.local            # Environment variables
├── Dockerfile            # Docker configuration
├── next.config.ts        # Next.js configuration
├── package.json          # NPM dependencies
└── tsconfig.json         # TypeScript configuration
```

## Docker Setup

The frontend can be run using Docker:

```bash
# Build the Docker image
docker build -t fullstack-frontend .

# Run the container
docker run -p 3000:3000 fullstack-frontend
```

For development with hot reloading, use Docker Compose from the root directory:

```bash
docker-compose up -d frontend
```

## Environment Variables

Create a `.env.local` file in the root directory with the following variables:

```
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

## Integration with Backend

This frontend application is designed to work with the Go backend API. The API client is configured to connect to the backend service using the `NEXT_PUBLIC_API_URL` environment variable.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.
