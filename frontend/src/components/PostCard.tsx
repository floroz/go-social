import React from 'react';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import type { Post } from "@/types/api"; // Assuming Post is re-exported from types/api.ts

interface PostCardProps {
  post: Post;
}

const PostCard: React.FC<PostCardProps> = ({ post }) => {
  // Helper to format date, can be moved to a utils file later
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <Card className="mb-4">
      <CardHeader>
        {/* TODO: Fetch and display author's username/name based on post.user_id */}
        <CardTitle className="text-lg">User ID: {post.user_id}</CardTitle>
        <p className="text-xs text-muted-foreground">
          Posted on: {formatDate(post.created_at)}
          {post.updated_at !== post.created_at && (
            <> | Updated: {formatDate(post.updated_at)}</>
          )}
        </p>
      </CardHeader>
      <CardContent>
        <p className="whitespace-pre-wrap">{post.content}</p>
      </CardContent>
      <CardFooter className="flex justify-between items-center">
        {/* TODO: Implement like functionality */}
        <Button variant="ghost" size="sm">
          {/* Icon for like, e.g., Heart */}
          Like (0)
        </Button>
        {/* TODO: Implement comment functionality */}
        <Button variant="ghost" size="sm">
          {/* Icon for comment, e.g., MessageSquare */}
          Comment (0)
        </Button>
        {/* TODO: Implement view comments functionality (e.g., navigate to post detail or expand) */}
        <Button variant="link" size="sm">
          View Details
        </Button>
      </CardFooter>
    </Card>
  );
};

export default PostCard;
