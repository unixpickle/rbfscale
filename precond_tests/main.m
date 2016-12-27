function [mat] = distance_mat(width, height, variance)
  mat = zeros(width*height, width*height);
  for x1 = 1:width
    for x2 = 1:width
      for y1 = 1:height
        for y2 = 1:height
          dist = (x1-x2)^2 + (y1-y2)^2;
          rbf = exp(-dist / (2*variance));
          mat(x1+(y1-1)*width, x2+(y2-1)*width) = rbf;
        end
      end
    end
  end
end

function [mat] = small_precond(width, height, variance)
  rbf1 = exp(-0.5 / variance);
  rbf2 = exp(-2 / variance);
  rbf3 = exp(-4.5 / variance);
  rbf4 = exp(-8 / variance);
  rbf5 = exp(-12.5 / variance);
  rbf6 = exp(-18 / variance);
  block = [1 rbf1 rbf2 rbf3 rbf4 rbf5 rbf6;
           rbf1 1 rbf1 rbf2 rbf3 rbf4 rbf5;
           rbf2 rbf1 1 rbf1 rbf2 rbf3 rbf4;
           rbf3 rbf2 rbf1 1 rbf1 rbf2 rbf3;
           rbf4 rbf3 rbf2 rbf1 1 rbf1 rbf2;
           rbf5 rbf4 rbf3 rbf2 rbf1 1 rbf1;
           rbf6 rbf5 rbf4 rbf3 rbf2 rbf1 1];
  block = block^(-1/2);
  n = width * height;
  mat = block(4,4) * eye(n);
  for i = 1:n
    if i > 1
      mat(i, i-1) = block(4, 3);
      if i > 2
        mat(i, i-2) = block(4, 2);
        if i > 3
          mat(i, i-3) = block(4, 1);
        end
      end
    end
    if i < n
      mat(i, i+1) = block(4, 5);
      if i < n-1
        mat(i, i+2) = block(4, 6);
        if i < n-2
          mat(i, i+3) = block(4, 7);
        end
      end
    end
  end
end

dm = distance_mat(10, 20, 3);
precond = small_precond(10, 20, 3);
cond(dm)
cond(precond*dm)
