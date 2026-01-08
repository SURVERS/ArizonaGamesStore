import React from 'react';

function StarRating({ rating }) {
  const maxStars = 5;
  const stars = [];

  for (let i = 1; i <= maxStars; i++) {
    const filled = rating >= i;
    const halfFilled = rating >= i - 0.5 && rating < i;

    stars.push(
      <span key={i} className="star">
        {filled ? '★' : halfFilled ? '⯨' : '☆'}
      </span>
    );
  }

  return (
    <div className="star-rating">
      {stars}
      <span className="rating-value">{rating.toFixed(1)}</span>
    </div>
  );
}

export default StarRating;
