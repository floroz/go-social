import React from 'react';
import { usePosts } from '@/hooks/usePosts'; // Corrected path assuming @ is src
import PostCard from '@/components/PostCard'; // Corrected path assuming @ is src
import { Button } from '@/components/ui/button';

const HomePage: React.FC = () => {
  const { posts, isLoading, error, refetchPosts } = usePosts();

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Home Feed</h1>
        <Button onClick={() => refetchPosts()} disabled={isLoading}>
          {isLoading ? 'Refreshing...' : 'Refresh Feed'}
        </Button>
      </div>

      {isLoading && <p>Loading posts...</p>}
      {error && (
        <div className="text-red-500">
          <p>Error fetching posts: {error.message}</p>
          <Button onClick={() => refetchPosts()} className="mt-2">
            Try Again
          </Button>
        </div>
      )}

      {!isLoading && !error && posts && posts.length === 0 && (
        <p>No posts yet. Be the first to create one!</p>
      )}

      {!isLoading && !error && posts && posts.length > 0 && (
        <div>
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </div>
      )}
      {/* TODO: Add Create Post functionality here or as a separate button/modal */}
    </div>
  );
};

export default HomePage;
