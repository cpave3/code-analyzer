using Shop.Domain;
using Xunit;

namespace Shop.Domain.Tests;

public sealed class OrderTests
{
    [Fact]
    public async Task Totals()
    {
        var order = new Order();
        Assert.Equal(3m, order.Total(new[] { 1m, 2m }));
        order.Total(new[] { 4m }).Should().Be(4m);
        await Assert.ThrowsAsync<OverflowException>(() => Task.Run(() => order.Total(new[] { decimal.MaxValue, decimal.MaxValue })));
    }
}
